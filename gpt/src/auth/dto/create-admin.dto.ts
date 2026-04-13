import {
  IsEmail,
  IsIn,
  IsOptional,
  IsString,
  Matches,
  MaxLength,
  MinLength,
} from 'class-validator';
import { ADMIN_USER_TYPES, ADMIN_VERIFIED_VALUES } from '../../common/constants';

export class CreateAdminDto {
  @IsString()
  @MaxLength(255)
  firstname!: string;

  @IsString()
  @MaxLength(255)
  lastname!: string;

  @IsOptional()
  @IsString()
  @MaxLength(32)
  @Matches(/^(?:\+?234|0)[789][01]\d{8}$|^(?:\+?1)?\d{10}$/, {
    message: 'phone must be a valid NG or US phone number',
  })
  phone?: string;

  @IsEmail()
  @MaxLength(255)
  email!: string;

  @IsIn(ADMIN_VERIFIED_VALUES)
  verified!: (typeof ADMIN_VERIFIED_VALUES)[number];

  @IsIn(ADMIN_USER_TYPES)
  userType!: (typeof ADMIN_USER_TYPES)[number];

  @IsString()
  @MinLength(6)
  @MaxLength(255)
  password!: string;
}
