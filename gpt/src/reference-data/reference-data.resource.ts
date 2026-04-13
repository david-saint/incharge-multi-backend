import { ContraceptionReason } from '../database/entities/contraception-reason.entity';
import { EducationLevel } from '../database/entities/education-level.entity';
import { FaqGroup } from '../database/entities/faq-group.entity';

export function buildContraceptionReasonResource(reason: ContraceptionReason) {
  return {
    id: reason.id,
    value: reason.value,
    profiles: reason.profiles?.map((profile) => ({ id: profile.id })),
  };
}

export function buildEducationLevelResource(level: EducationLevel) {
  return {
    id: level.id,
    name: level.name,
  };
}

export function buildFaqGroupResource(group: FaqGroup) {
  return {
    id: group.id,
    name: group.name,
    created_at: group.createdAt,
    updated_at: group.updatedAt,
  };
}
