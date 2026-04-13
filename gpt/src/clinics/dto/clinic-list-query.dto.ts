import { IsOptional, IsString } from 'class-validator';

export class ClinicListQueryDto {
  @IsOptional()
  @IsString()
  search?: string;

  @IsOptional()
  @IsString()
  latitude?: string;

  @IsOptional()
  @IsString()
  longitude?: string;

  @IsOptional()
  @IsString()
  radius?: string;

  @IsOptional()
  @IsString()
  mode?: string;

  @IsOptional()
  @IsString()
  sort?: string;

  @IsOptional()
  @IsString()
  with?: string;

  @IsOptional()
  @IsString()
  page?: string;

  @IsOptional()
  @IsString()
  per_page?: string;

  @IsOptional()
  @IsString()
  withTrashed?: string;

  @IsOptional()
  @IsString()
  onlyTrashed?: string;

  @IsOptional()
  @IsString()
  withCount?: string;
}
