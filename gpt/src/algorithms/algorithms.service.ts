import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { IsNull, Repository } from 'typeorm';
import { Algorithm } from '../database/entities/algorithm.entity';
import { buildAlgorithmResource } from './algorithm.resource';
import { SaveAlgorithmDto } from './dto/save-algorithm.dto';
import { UpdateAlgorithmDto } from './dto/update-algorithm.dto';

@Injectable()
export class AlgorithmsService {
  constructor(
    @InjectRepository(Algorithm)
    private readonly algorithmRepository: Repository<Algorithm>,
  ) {}

  async list(active?: string) {
    const where:
      | { active: 'Y' | 'N'; deletedAt: ReturnType<typeof IsNull> }
      | { deletedAt: ReturnType<typeof IsNull> } =
      active === 'Y' || active === 'N'
        ? { active, deletedAt: IsNull() }
        : { deletedAt: IsNull() };
    const algorithms = await this.algorithmRepository.find({
      where,
      order: { active: 'ASC', id: 'ASC' },
    });
    return algorithms.map(buildAlgorithmResource);
  }

  async create(payload: SaveAlgorithmDto) {
    const algorithm = this.algorithmRepository.create({
      ...payload,
      active: payload.active ?? 'N',
    });
    const saved = await this.algorithmRepository.save(algorithm);
    return buildAlgorithmResource(saved);
  }

  async update(id: number, payload: UpdateAlgorithmDto) {
    const algorithm = await this.algorithmRepository.findOne({ where: { id, deletedAt: IsNull() } });
    if (!algorithm) {
      throw new NotFoundException();
    }
    Object.assign(algorithm, payload);
    const saved = await this.algorithmRepository.save(algorithm);
    return buildAlgorithmResource(saved);
  }
}
