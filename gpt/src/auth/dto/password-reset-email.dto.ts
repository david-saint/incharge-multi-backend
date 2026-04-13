import { IsEmail, MaxLength } from 'class-validator';

export class PasswordResetEmailDto {
  @IsEmail()
  @MaxLength(255)
  email!: string;
}
