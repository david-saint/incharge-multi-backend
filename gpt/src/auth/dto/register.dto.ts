import {
  IsEmail,
  IsOptional,
  IsString,
  Matches,
  MaxLength,
  MinLength,
} from 'class-validator';

export class RegisterDto {
  @IsString()
  @MaxLength(255)
  name!: string;

  @IsEmail()
  @MaxLength(255)
  email!: string;

  @IsOptional()
  @IsString()
  @MaxLength(32)
  @Matches(/^(?:\+?234|0)[789][01]\d{8}$|^(?:\+?1)?\d{10}$/, {
    message: 'phone must be a valid NG or US phone number',
  })
  phone?: string;

  @IsString()
  @MinLength(6)
  @MaxLength(255)
  password!: string;
}
