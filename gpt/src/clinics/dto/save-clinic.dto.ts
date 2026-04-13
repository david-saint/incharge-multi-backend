import { IsNumber, IsString, MaxLength } from 'class-validator';

export class SaveClinicDto {
  @IsString()
  @MaxLength(255)
  name!: string;

  @IsString()
  address!: string;

  @IsNumber()
  latitude!: number;

  @IsNumber()
  longitude!: number;

  @IsNumber()
  added_by_id!: number;
}
