import { PROFILE_PLAN_KEY } from '../common/constants';
import { Profile } from '../database/entities/profile.entity';
import {
  buildContraceptionReasonResource,
  buildEducationLevelResource,
} from '../reference-data/reference-data.resource';

export interface ProfileResource {
  id: number;
  age: number;
  gender: string;
  date_of_birth: Date;
  address: string;
  latitude: string | null;
  longitude: string | null;
  marital_status: string;
  height: number | null;
  weight: number | null;
  occupation: string | null;
  children: number;
  sexually_active: boolean;
  pregnancy_status: boolean;
  religion: string;
  religion_sect: string | null;
  contraceptive_plan?: unknown;
  reason?: unknown;
  education_level?: unknown;
  user?: {
    id: number;
    name: string;
    email: string;
    phone: string | null;
    email_verified: boolean;
    phone_confirmed: boolean;
    created_at: Date;
    updated_at: Date;
  };
}

export function buildProfileResource(
  profile: Profile,
  options?: { includeUser?: boolean },
): ProfileResource {
  return {
    id: profile.id,
    age: profile.age,
    gender: profile.gender,
    date_of_birth: profile.dateOfBirth,
    address: profile.address,
    latitude: profile.latitude,
    longitude: profile.longitude,
    marital_status: profile.maritalStatus,
    height: profile.height,
    weight: profile.weight === null ? null : Number(profile.weight),
    occupation: profile.occupation,
    children: profile.numberOfChildren ?? 0,
    sexually_active: profile.sexuallyActive,
    pregnancy_status: profile.pregnancyStatus,
    religion: profile.religion,
    religion_sect: profile.religionSect,
    contraceptive_plan:
      profile.meta && PROFILE_PLAN_KEY in profile.meta
        ? profile.meta[PROFILE_PLAN_KEY]
        : undefined,
    reason: profile.contraceptionReason
      ? buildContraceptionReasonResource(profile.contraceptionReason)
      : undefined,
    education_level: profile.educationLevel
      ? buildEducationLevelResource(profile.educationLevel)
      : undefined,
    user:
      options?.includeUser && profile.user
        ? {
            id: profile.user.id,
            name: profile.user.name,
            email: profile.user.email,
            phone: profile.user.phone,
            email_verified: Boolean(profile.user.emailVerifiedAt),
            phone_confirmed: Boolean(profile.user.phoneConfirmedAt),
            created_at: profile.user.createdAt,
            updated_at: profile.user.updatedAt,
          }
        : undefined,
  };
}
