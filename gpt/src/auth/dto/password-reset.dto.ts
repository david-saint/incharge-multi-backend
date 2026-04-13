import { IsEmail, IsString, MaxLength, MinLength } from 'class-validator';

export class PasswordResetDto {
  @IsEmail()
  @MaxLength(255)
  email!: string;

  @IsString()
  @MaxLength(255)
  token!: string;

  @IsString()
  @MinLength(6)
  @MaxLength(255)
  password!: string;

  @IsString()
  @MinLength(6)
  @MaxLength(255)
  password_confirmation!: string;
}
