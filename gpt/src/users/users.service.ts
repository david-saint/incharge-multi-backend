import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { IsNull, Not, Repository } from 'typeorm';
import { buildPaginationMeta } from '../common/utils/query.util';
import { User } from '../database/entities/user.entity';
import { buildUserResource } from './user.resource';

@Injectable()
export class UsersService {
  constructor(
    @InjectRepository(User)
    private readonly userRepository: Repository<User>,
  ) {}

  async listDeleted(includeDeleted = false, basePath?: string) {
    const users = await this.userRepository.find({
      where: includeDeleted ? { deletedAt: Not(IsNull()) } : { deletedAt: IsNull() },
      withDeleted: includeDeleted,
      relations: {
        profile: {
          educationLevel: true,
          contraceptionReason: true,
        },
      },
      order: { id: 'ASC' },
      take: 50,
    });

    return {
      data: users.map((user) => buildUserResource(user, { includeProfile: true })),
      ...buildPaginationMeta(
        basePath ?? (includeDeleted ? '/getDeletedUsers' : '/getUsers'),
        1,
        50,
        users.length,
      ),
    };
  }

  async softDelete(id: number) {
    const user = await this.userRepository.findOne({ where: { id } });
    if (!user) {
      throw new NotFoundException();
    }
    await this.userRepository.softDelete(id);
    return { status: true, message: 'User deleted successfully.' };
  }

  async restore(id: number) {
    const result = await this.userRepository.restore(id);
    if (!result.affected) {
      throw new NotFoundException();
    }
    return { status: true, message: 'User restored successfully.' };
  }
}
