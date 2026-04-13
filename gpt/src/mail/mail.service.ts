import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import nodemailer from 'nodemailer';
import { User } from '../database/entities/user.entity';
import { signPayload } from '../common/utils/security.util';

@Injectable()
export class MailService {
  private readonly logger = new Logger(MailService.name);

  constructor(private readonly config: ConfigService) {}

  async sendVerificationEmail(user: User): Promise<void> {
    const verificationLink = this.buildVerificationLink(user);
    await this.sendMail(
      user.email,
      'Verify your email',
      `Welcome to InCharge. Verify your email using this link: ${verificationLink}`,
    );
  }

  async sendPasswordResetEmail(user: User, token: string): Promise<void> {
    const url = `${this.config.get<string>('app.userDomain', 'http://localhost:5173')}/reset-password/${token}`;
    await this.sendMail(
      user.email,
      'Reset your password',
      `Reset your InCharge password using this link: ${url}`,
    );
  }

  private buildVerificationLink(user: User): string {
    const expires = Date.now() +
      this.config.get<number>('app.emailVerificationExpiryMinutes', 4320) * 60_000;
    const signature = signPayload(
      `${user.id}:${expires}`,
      this.config.get<string>('app.jwtSecret', 'change-me'),
    );
    return `${this.config.get<string>('app.apiDomain', 'http://localhost:3000')}/api/v1/user/email/verify/${user.id}?expires=${expires}&signature=${signature}`;
  }

  private async sendMail(to: string, subject: string, text: string): Promise<void> {
    if (!this.config.get<boolean>('app.mailEnabled', false)) {
      this.logger.log(`Mail disabled; would send to ${to}: ${subject}`);
      return;
    }

    const transporter = nodemailer.createTransport({
      host: this.config.get<string>('MAIL_HOST', 'localhost'),
      port: this.config.get<number>('MAIL_PORT', 1025),
      secure: ['ssl', 'tls'].includes(
        this.config.get<string>('MAIL_ENCRYPTION', '').toLowerCase(),
      ),
      auth:
        this.config.get<string>('MAIL_USERNAME')
          ? {
              user: this.config.get<string>('MAIL_USERNAME'),
              pass: this.config.get<string>('MAIL_PASSWORD'),
            }
          : undefined,
    });

    await transporter.sendMail({
      from: 'InCharge <no-reply@incharge.local>',
      to,
      subject,
      text,
    });
  }
}
