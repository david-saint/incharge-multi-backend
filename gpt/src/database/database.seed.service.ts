import { existsSync } from 'node:fs';
import { readFile } from 'node:fs/promises';
import { join } from 'node:path';
import { Injectable, Logger, OnApplicationBootstrap } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { InjectRepository } from '@nestjs/typeorm';
import { parse } from 'csv-parse/sync';
import { DataSource, IsNull, Repository } from 'typeorm';
import { Algorithm } from './entities/algorithm.entity';
import { ContraceptionReason } from './entities/contraception-reason.entity';
import { Country } from './entities/country.entity';
import { EducationLevel } from './entities/education-level.entity';
import { Faq } from './entities/faq.entity';
import { FaqGroup } from './entities/faq-group.entity';
import { Location } from './entities/location.entity';
import { State } from './entities/state.entity';

@Injectable()
export class DatabaseSeedService implements OnApplicationBootstrap {
  private readonly logger = new Logger(DatabaseSeedService.name);

  constructor(
    private readonly config: ConfigService,
    private readonly dataSource: DataSource,
    @InjectRepository(ContraceptionReason)
    private readonly contraceptionReasonRepository: Repository<ContraceptionReason>,
    @InjectRepository(EducationLevel)
    private readonly educationLevelRepository: Repository<EducationLevel>,
    @InjectRepository(FaqGroup)
    private readonly faqGroupRepository: Repository<FaqGroup>,
    @InjectRepository(Faq)
    private readonly faqRepository: Repository<Faq>,
    @InjectRepository(Country)
    private readonly countryRepository: Repository<Country>,
    @InjectRepository(State)
    private readonly stateRepository: Repository<State>,
    @InjectRepository(Location)
    private readonly locationRepository: Repository<Location>,
    @InjectRepository(Algorithm)
    private readonly algorithmRepository: Repository<Algorithm>,
  ) {}

  async onApplicationBootstrap(): Promise<void> {
    if (!this.config.get<boolean>('app.autoSeed', false)) {
      return;
    }
    await this.seedReferenceData();
  }

  async seedReferenceData(): Promise<void> {
    await this.seedContraceptionReasons();
    await this.seedEducationLevels();
    await this.seedFaqGroups();
    await this.seedGeoData();
    await this.seedAlgorithms();
  }

  private async seedContraceptionReasons(): Promise<void> {
    const count = await this.contraceptionReasonRepository.count({
      where: { deletedAt: IsNull() },
    });
    if (count > 0) {
      return;
    }

    await this.contraceptionReasonRepository.save(
      [
        'Completed family size',
        'Child Spacing',
        'Sexually Active with no plan for children at the moment',
      ].map((value) => this.contraceptionReasonRepository.create({ value })),
    );
  }

  private async seedEducationLevels(): Promise<void> {
    const count = await this.educationLevelRepository.count();
    if (count > 0) {
      return;
    }

    const values = [
      'BArch',
      'BA',
      'B.Sc',
      'B.ENG',
      'LLB',
      'HNC',
      'HND',
      'ND',
      'M.Sc',
      'M.Eng',
      'Phd.',
      'Prof',
      'B.Tech',
      'Other',
    ];
    await this.educationLevelRepository.save(
      values.map((name) => this.educationLevelRepository.create({ name })),
    );
  }

  private async seedFaqGroups(): Promise<void> {
    const count = await this.faqGroupRepository.count();
    if (count > 0) {
      return;
    }

    const groups = [
      'Barrier Method',
      'Combined Oral Contraceptives',
      'Diaphragms and Spermicides',
      'Emergency Contraceptive Pills',
      'Female Sterilization',
      'Fertility Awareness',
      'Implants',
      'Injectables',
      'IUCD',
      'Lactational Amenorrhea',
      'Progestin Only Pills',
      'STIs',
      'Vasectomy',
    ];

    const createdGroups = await this.faqGroupRepository.save(
      groups.map((name) => this.faqGroupRepository.create({ name })),
    );

    await this.faqRepository.save(
      createdGroups.map((group) =>
        this.faqRepository.create({
          faqGroupId: group.id,
          content: {
            title: group.name,
            blocks: [
              {
                type: 'paragraph',
                text: `Reference content for ${group.name}`,
              },
            ],
          },
        }),
      ),
    );
  }

  private async seedGeoData(): Promise<void> {
    const countryCount = await this.countryRepository.count();
    if (countryCount > 0) {
      return;
    }

    const countriesPath = join(process.cwd(), 'database', 'data', 'countries.csv');
    const statesPath = join(process.cwd(), 'database', 'data', 'states.csv');
    const locationsPath = join(process.cwd(), 'database', 'data', 'locations.csv');

    const countries = existsSync(countriesPath)
      ? await this.loadCsv(countriesPath)
      : [
          { id: '1', name: 'Nigeria', code: 'NG' },
          { id: '2', name: 'United States', code: 'US' },
        ];
    const states = existsSync(statesPath)
      ? await this.loadCsv(statesPath)
      : [
          {
            id: '1',
            name: 'Lagos',
            slug: 'lagos',
            latitude: '6.5244',
            longitude: '3.3792',
          },
          {
            id: '2',
            name: 'Abuja',
            slug: 'abuja',
            latitude: '9.0765',
            longitude: '7.3986',
          },
        ];
    const locations = existsSync(locationsPath)
      ? await this.loadCsv(locationsPath)
      : [
          {
            id: '1',
            name: 'Ikeja',
            state_id: '1',
            country_id: '1',
            latitude: '6.6018',
            longitude: '3.3515',
          },
          {
            id: '2',
            name: 'Maitama',
            state_id: '2',
            country_id: '1',
            latitude: '9.0822',
            longitude: '7.4951',
          },
        ];

    await this.countryRepository.save(
      countries.map((country) =>
        this.countryRepository.create({
          id: Number(country.id),
          name: country.name,
          code: country.code,
        }),
      ),
    );

    await this.stateRepository.save(
      states.map((state) =>
        this.stateRepository.create({
          id: Number(state.id),
          name: state.name,
          slug: state.slug,
          latitude: state.latitude ?? null,
          longitude: state.longitude ?? null,
          meta: null,
        }),
      ),
    );

    await this.locationRepository.save(
      locations.map((location) =>
        this.locationRepository.create({
          id: Number(location.id),
          name: location.name,
          stateId: Number(location.state_id),
          countryId: Number(location.country_id),
          latitude: location.latitude ?? null,
          longitude: location.longitude ?? null,
          meta: null,
        }),
      ),
    );
  }

  private async seedAlgorithms(): Promise<void> {
    const count = await this.algorithmRepository.count();
    if (count > 0) {
      return;
    }

    const seedPath = join(process.cwd(), 'database', 'sql', 'algorithms.sql');
    if (!existsSync(seedPath)) {
      this.logger.warn('No algorithms.sql file found; skipping algorithm seed');
      return;
    }

    const sql = await readFile(seedPath, 'utf8');
    if (!sql.trim()) {
      return;
    }
    await this.dataSource.query(sql);
  }

  private async loadCsv(path: string): Promise<Record<string, string>[]> {
    const file = await readFile(path, 'utf8');
    return parse(file, { columns: true, skip_empty_lines: true });
  }
}
