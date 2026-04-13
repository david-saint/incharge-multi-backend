import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { IsNull, Repository } from 'typeorm';
import { ContraceptionReason } from '../database/entities/contraception-reason.entity';
import { EducationLevel } from '../database/entities/education-level.entity';
import { Faq } from '../database/entities/faq.entity';
import { FaqGroup } from '../database/entities/faq-group.entity';
import {
  buildContraceptionReasonResource,
  buildEducationLevelResource,
  buildFaqGroupResource,
} from '../reference-data/reference-data.resource';

@Injectable()
export class GlobalService {
  constructor(
    @InjectRepository(ContraceptionReason)
    private readonly reasonRepository: Repository<ContraceptionReason>,
    @InjectRepository(EducationLevel)
    private readonly educationLevelRepository: Repository<EducationLevel>,
    @InjectRepository(FaqGroup)
    private readonly faqGroupRepository: Repository<FaqGroup>,
    @InjectRepository(Faq)
    private readonly faqRepository: Repository<Faq>,
  ) {}

  hello(): string {
    return 'Hello, World!';
  }

  async listContraceptionReasons() {
    const reasons = await this.reasonRepository.find({
      where: { deletedAt: IsNull() },
      order: { id: 'ASC' },
    });
    return reasons.map(buildContraceptionReasonResource);
  }

  async getContraceptionReason(id: number) {
    const reason = await this.reasonRepository.findOne({
      where: { id, deletedAt: IsNull() },
    });
    if (!reason) {
      throw new NotFoundException();
    }
    return buildContraceptionReasonResource(reason);
  }

  async listEducationLevels() {
    const levels = await this.educationLevelRepository.find({ order: { id: 'ASC' } });
    return levels.map(buildEducationLevelResource);
  }

  async listFaqGroups() {
    const groups = await this.faqGroupRepository.find({ order: { id: 'ASC' } });
    return groups.map(buildFaqGroupResource);
  }

  async getFaqGroupContent(id: number) {
    const faq = await this.faqRepository.findOne({ where: { faqGroupId: id, deletedAt: IsNull() } });
    if (!faq) {
      throw new NotFoundException();
    }
    return {
      data: faq.content,
      status: 'faq.get_content',
    };
  }
}
