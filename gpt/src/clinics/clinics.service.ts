import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { In, IsNull, Like, Not, Repository } from 'typeorm';
import { CLINIC_LOCATABLE_TYPE } from '../common/constants';
import { haversineDistance } from '../common/utils/geo.util';
import {
  buildPaginationMeta,
  parseBooleanQuery,
  parseCsvList,
  parseSortDirectives,
  toNumberOrNull,
} from '../common/utils/query.util';
import { Clinic } from '../database/entities/clinic.entity';
import { Locatable } from '../database/entities/locatable.entity';
import { Location } from '../database/entities/location.entity';
import { buildClinicResource } from './clinic.resource';
import { ClinicListQueryDto } from './dto/clinic-list-query.dto';
import { SaveClinicDto } from './dto/save-clinic.dto';

@Injectable()
export class ClinicsService {
  constructor(
    @InjectRepository(Clinic)
    private readonly clinicRepository: Repository<Clinic>,
    @InjectRepository(Locatable)
    private readonly locatableRepository: Repository<Locatable>,
    @InjectRepository(Location)
    private readonly locationRepository: Repository<Location>,
  ) {}

  async list(query: ClinicListQueryDto, basePath: string) {
    const includeLocations = parseCsvList(query.with).includes('locations');
    const includeLocationsCount = parseCsvList(query.withCount).includes('locations');
    const withTrashed = parseBooleanQuery(query.withTrashed);
    const onlyTrashed = parseBooleanQuery(query.onlyTrashed);
    const page = query.page
      ? Math.max(1, Number.parseInt(query.page, 10) || 1)
      : null;
    const perPage = Math.max(1, Number.parseInt(query.per_page ?? '20', 10));
    const latitude = toNumberOrNull(query.latitude);
    const longitude = toNumberOrNull(query.longitude);
    const distanceMode = query.mode === 'mi' ? 'mi' : 'km';
    const radius = toNumberOrNull(query.radius) ?? 10;
    const hasDistanceFilter = latitude !== null && longitude !== null;

    const clinics = await this.clinicRepository.find({
      where: onlyTrashed
        ? { deletedAt: Not(IsNull()) }
        : withTrashed
          ? {}
          : { deletedAt: IsNull() },
      withDeleted: withTrashed || onlyTrashed,
      relations: includeLocations
        ? {
            locatables: {
              location: {
                state: true,
                country: true,
              },
            },
          }
        : undefined,
      order: { id: 'ASC' },
    });

    const locationMatches = query.search
      ? await this.locationRepository.find({
          where: { name: Like(`%${query.search}%`) },
        })
      : [];
    const matchedLocatables = locationMatches.length
      ? await this.locatableRepository.find({
          where: {
            locationId: In(locationMatches.map((location) => location.id)),
            locatableType: CLINIC_LOCATABLE_TYPE,
          },
        })
      : [];
    const locationClinicIds = new Set(matchedLocatables.map((item) => item.locatableId));

    const filtered = clinics
      .filter((clinic) => {
        if (!query.search) {
          return true;
        }

        const needle = query.search.toLowerCase();
        return (
          String(clinic.id).includes(needle) ||
          clinic.name.toLowerCase().includes(needle) ||
          clinic.address.toLowerCase().includes(needle) ||
          locationClinicIds.has(clinic.id)
        );
      })
      .map((clinic) => {
        const distance =
          hasDistanceFilter && clinic.latitude && clinic.longitude
            ? haversineDistance(
                latitude!,
                longitude!,
                Number(clinic.latitude),
                Number(clinic.longitude),
                distanceMode,
              )
            : null;
        return { clinic, distance };
      })
      .filter((entry) => !hasDistanceFilter || (entry.distance !== null && entry.distance <= radius));

    const locationCountMap =
      includeLocationsCount && filtered.length > 0
        ? await this.buildLocationCountMap(filtered.map((entry) => entry.clinic.id))
        : new Map<number, number>();

    const sortDirectives = parseSortDirectives(query.sort, [
      'id',
      'name',
      'address',
      'created_at',
      'distance',
    ]);

    filtered.sort((left, right) => {
      for (const directive of sortDirectives) {
        const direction = directive.direction === 'ASC' ? 1 : -1;
        const leftValue = this.sortValue(left.clinic, left.distance, directive.column);
        const rightValue = this.sortValue(right.clinic, right.distance, directive.column);
        if (leftValue < rightValue) {
          return -1 * direction;
        }
        if (leftValue > rightValue) {
          return 1 * direction;
        }
      }
      return 0;
    });

    const total = filtered.length;
    if (page === null) {
      return filtered.map((entry) =>
        buildClinicResource(entry.clinic, {
          includeLocations,
          locationsCount: includeLocationsCount
            ? locationCountMap.get(entry.clinic.id) ?? 0
            : undefined,
          distance:
            entry.distance === null
              ? undefined
              : { mode: distanceMode, radius, actualDistance: entry.distance },
        }),
      );
    }

    const start = (page - 1) * perPage;
    const data = filtered.slice(start, start + perPage).map((entry) =>
      buildClinicResource(entry.clinic, {
        includeLocations,
        locationsCount: includeLocationsCount
          ? locationCountMap.get(entry.clinic.id) ?? 0
          : undefined,
        distance:
          entry.distance === null
            ? undefined
            : { mode: distanceMode, radius, actualDistance: entry.distance },
      }),
    );

    return {
      data,
      ...buildPaginationMeta(basePath, page, perPage, total),
    };
  }

  async create(payload: SaveClinicDto) {
    const clinic = this.clinicRepository.create({
      name: payload.name,
      address: payload.address,
      latitude: payload.latitude.toFixed(7),
      longitude: payload.longitude.toFixed(7),
      addedById: payload.added_by_id,
      meta: null,
    });

    const saved = await this.clinicRepository.save(clinic);
    return buildClinicResource(saved);
  }

  async update(id: number, payload: SaveClinicDto) {
    const clinic = await this.clinicRepository.findOne({ where: { id } });
    if (!clinic) {
      throw new NotFoundException();
    }
    clinic.name = payload.name;
    clinic.address = payload.address;
    clinic.latitude = payload.latitude.toFixed(7);
    clinic.longitude = payload.longitude.toFixed(7);
    clinic.addedById = payload.added_by_id;
    const saved = await this.clinicRepository.save(clinic);
    return buildClinicResource(saved);
  }

  async softDelete(id: number) {
    const clinic = await this.clinicRepository.findOne({ where: { id } });
    if (!clinic) {
      throw new NotFoundException();
    }
    await this.clinicRepository.softDelete(id);
    return { status: true, message: 'Clinic deleted successfully.' };
  }

  async restore(id: number) {
    const result = await this.clinicRepository.restore(id);
    if (!result.affected) {
      throw new NotFoundException();
    }
    return { status: true, message: 'Clinic restored successfully.' };
  }

  async listSimple(includeDeleted = false, basePath?: string) {
    return this.list(
      { page: '1', per_page: '50', onlyTrashed: includeDeleted ? 'true' : undefined },
      basePath ?? (includeDeleted ? '/getDeletedClinics' : '/getClinics'),
    );
  }

  private async buildLocationCountMap(clinicIds: number[]) {
    const locatables = await this.locatableRepository.find({
      where: {
        locatableId: In(clinicIds),
        locatableType: CLINIC_LOCATABLE_TYPE,
      },
    });

    return locatables.reduce<Map<number, number>>((map, locatable) => {
      map.set(locatable.locatableId, (map.get(locatable.locatableId) ?? 0) + 1);
      return map;
    }, new Map<number, number>());
  }

  private sortValue(clinic: Clinic, distance: number | null, column: string) {
    switch (column) {
      case 'name':
        return clinic.name.toLowerCase();
      case 'address':
        return clinic.address.toLowerCase();
      case 'created_at':
        return clinic.createdAt.getTime();
      case 'distance':
        return distance ?? Number.MAX_SAFE_INTEGER;
      case 'id':
      default:
        return clinic.id;
    }
  }
}
