import { IsString, MaxLength } from 'class-validator';

export class StorePlanDto {
  @IsString()
  @MaxLength(255)
  plan!: string;
}
