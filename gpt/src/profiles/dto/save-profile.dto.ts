import {
  IsBoolean,
  IsDateString,
  IsIn,
  IsNumber,
  IsOptional,
  IsString,
  MaxLength,
  ValidateIf,
} from 'class-validator';
import {
  GENDERS,
  MARITAL_STATUSES,
  RELIGIONS,
  RELIGION_SECTS,
} from '../../common/constants';

export class SaveProfileDto {
  @IsIn(GENDERS)
  gender!: (typeof GENDERS)[number];

  @IsOptional()
  @IsNumber()
  age?: number;

  @IsOptional()
  @IsDateString()
  dob?: string;

  @IsOptional()
  @IsString()
  address?: string;

  @IsOptional()
  @IsIn(MARITAL_STATUSES)
  marital_status?: (typeof MARITAL_STATUSES)[number];

  @IsOptional()
  @IsNumber()
  height?: number;

  @IsOptional()
  @IsNumber()
  weight?: number;

  @IsOptional()
  @IsNumber()
  education_level?: number;

  @IsOptional()
  @IsString()
  @MaxLength(255)
  occupation?: string;

  @IsOptional()
  @IsNumber()
  children?: number;

  @IsOptional()
  @IsNumber()
  reason?: number;

  @IsOptional()
  @IsBoolean()
  sexually_active?: boolean;

  @IsOptional()
  @IsBoolean()
  pregnancy_status?: boolean;

  @IsOptional()
  @IsIn(RELIGIONS)
  religion?: (typeof RELIGIONS)[number];

  @ValidateIf((value: SaveProfileDto) => value.religion === 'CHRISTIANITY')
  @IsIn(RELIGION_SECTS)
  religion_sect?: (typeof RELIGION_SECTS)[number];
}
