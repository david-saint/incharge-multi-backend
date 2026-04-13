import { Injectable, NotFoundException, UnprocessableEntityException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { PROFILE_PLAN_KEY } from '../common/constants';
import { EducationLevel } from '../database/entities/education-level.entity';
import { Profile } from '../database/entities/profile.entity';
import { User } from '../database/entities/user.entity';
import { ContraceptionReason } from '../database/entities/contraception-reason.entity';
import { buildProfileResource } from './profile.resource';
import { SaveProfileDto } from './dto/save-profile.dto';

@Injectable()
export class ProfilesService {
  constructor(
    @InjectRepository(Profile)
    private readonly profileRepository: Repository<Profile>,
    @InjectRepository(User)
    private readonly userRepository: Repository<User>,
    @InjectRepository(EducationLevel)
    private readonly educationLevelRepository: Repository<EducationLevel>,
    @InjectRepository(ContraceptionReason)
    private readonly reasonRepository: Repository<ContraceptionReason>,
  ) {}

  async upsertProfile(userId: number, payload: SaveProfileDto) {
    const user = await this.userRepository.findOne({ where: { id: userId } });
    if (!user) {
      throw new NotFoundException();
    }

    const educationLevelId = payload.education_level ?? 14;
    const reasonId = payload.reason ?? 3;

    await this.assertReferenceData(educationLevelId, reasonId);

    const existing = await this.profileRepository.findOne({ where: { userId } });
    const profile = existing ?? this.profileRepository.create({ userId, meta: {} });
    profile.age = payload.age ?? 0;
    profile.gender = payload.gender;
    profile.dateOfBirth = payload.dob ? new Date(payload.dob) : new Date();
    profile.address = payload.address ?? '';
    profile.maritalStatus = payload.marital_status ?? 'SINGLE';
    profile.height = payload.height ?? null;
    profile.weight = payload.weight === undefined ? null : payload.weight.toFixed(2);
    profile.educationLevelId = educationLevelId;
    profile.occupation = payload.occupation ?? null;
    profile.numberOfChildren = payload.children ?? 0;
    profile.contraceptionReasonId = reasonId;
    profile.sexuallyActive = payload.sexually_active ?? false;
    profile.pregnancyStatus = payload.pregnancy_status ?? false;
    profile.religion = payload.religion ?? 'OTHER';
    profile.religionSect =
      profile.religion === 'CHRISTIANITY' ? payload.religion_sect ?? null : null;

    const saved = await this.profileRepository.save(profile);
    const hydrated = await this.getProfileEntity(userId, ['reason', 'educationLevel']);
    return buildProfileResource(hydrated ?? saved);
  }

  async getProfile(userId: number, withRelations: string[]) {
    const profile = await this.getProfileEntity(userId, withRelations);
    if (!profile) {
      throw new NotFoundException();
    }

    return buildProfileResource(profile, {
      includeUser: withRelations.includes('user'),
    });
  }

  async storePlan(userId: number, plan: string) {
    const profile = await this.profileRepository.findOne({ where: { userId } });
    if (!profile) {
      throw new NotFoundException({ message: 'Profile not found.' });
    }

    profile.meta = { ...(profile.meta ?? {}), [PROFILE_PLAN_KEY]: plan };
    await this.profileRepository.save(profile);
    return {
      status: true,
      message: 'Contraceptive plan stored successfully.',
      data: plan,
    };
  }

  private async getProfileEntity(userId: number, withRelations: string[]) {
    const includeUser = withRelations.includes('user');
    const includeReason = withRelations.includes('reason');
    const includeEducationLevel = withRelations.includes('educationLevel');

    return this.profileRepository.findOne({
      where: { userId },
      relations: {
        user: includeUser,
        contraceptionReason: includeReason,
        educationLevel: includeEducationLevel,
      },
    });
  }

  private async assertReferenceData(educationLevelId: number, reasonId: number) {
    const [educationLevel, reason] = await Promise.all([
      this.educationLevelRepository.findOne({ where: { id: educationLevelId } }),
      this.reasonRepository.findOne({ where: { id: reasonId } }),
    ]);

    if (!educationLevel) {
      throw new UnprocessableEntityException({
        errors: { education_level: ['The selected education level is invalid.'] },
      });
    }
    if (!reason) {
      throw new UnprocessableEntityException({
        errors: { reason: ['The selected contraception reason is invalid.'] },
      });
    }
  }
}
