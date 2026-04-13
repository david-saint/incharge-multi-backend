import { Profile } from '../database/entities/profile.entity';
import { User } from '../database/entities/user.entity';
import {
  buildProfileResource,
  type ProfileResource,
} from '../profiles/profile.resource';

export interface UserResource {
  id: number;
  name: string;
  email: string;
  phone: string | null;
  email_verified: boolean;
  phone_confirmed: boolean;
  profile?: ProfileResource;
  created_at: Date;
  updated_at: Date;
}

export function buildUserResource(
  user: User,
  options?: { includeProfile?: boolean },
): UserResource {
  return {
    id: user.id,
    name: user.name,
    email: user.email,
    phone: user.phone,
    email_verified: Boolean(user.emailVerifiedAt),
    phone_confirmed: Boolean(user.phoneConfirmedAt),
    profile:
      options?.includeProfile && user.profile
        ? buildProfileResource(user.profile as Profile)
        : undefined,
    created_at: user.createdAt,
    updated_at: user.updatedAt,
  };
}
