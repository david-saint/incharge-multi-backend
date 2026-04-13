import { randomUUID } from 'node:crypto';
import {
  Injectable,
  NotFoundException,
  UnauthorizedException,
  UnprocessableEntityException,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { JwtService } from '@nestjs/jwt';
import { InjectRepository } from '@nestjs/typeorm';
import * as bcrypt from 'bcrypt';
import type { Request } from 'express';
import { IsNull, Repository } from 'typeorm';
import { Admin } from '../database/entities/admin.entity';
import { PasswordReset } from '../database/entities/password-reset.entity';
import { UserJwtSession } from '../database/entities/user-jwt-session.entity';
import { User } from '../database/entities/user.entity';
import { MailService } from '../mail/mail.service';
import { buildUserResource } from '../users/user.resource';
import {
  constantTimeEquals,
  generateOpaqueToken,
  sha256,
  signPayload,
} from '../common/utils/security.util';
import { successResponse } from '../common/utils/response.util';
import { LoginDto } from './dto/login.dto';
import { PasswordResetDto } from './dto/password-reset.dto';
import { PasswordResetEmailDto } from './dto/password-reset-email.dto';
import { RegisterDto } from './dto/register.dto';
import type { JwtPayload } from './jwt-payload.interface';

type AuthenticatedUser = User & { jwtJti?: string };

@Injectable()
export class AuthService {
  constructor(
    private readonly config: ConfigService,
    private readonly jwtService: JwtService,
    private readonly mailService: MailService,
    @InjectRepository(User)
    private readonly userRepository: Repository<User>,
    @InjectRepository(UserJwtSession)
    private readonly userJwtSessionRepository: Repository<UserJwtSession>,
    @InjectRepository(PasswordReset)
    private readonly passwordResetRepository: Repository<PasswordReset>,
    @InjectRepository(Admin)
    private readonly adminRepository: Repository<Admin>,
  ) {}

  async register(payload: RegisterDto) {
    await this.ensureUniqueUser(payload.email, payload.phone ?? null);
    const user = this.userRepository.create({
      name: payload.name,
      email: payload.email.toLowerCase(),
      phone: payload.phone ?? null,
      password: await bcrypt.hash(payload.password, 10),
      emailVerifiedAt: null,
      phoneConfirmedAt: null,
      rememberToken: null,
    });
    const savedUser = await this.userRepository.save(user);
    await this.mailService.sendVerificationEmail(savedUser);

    return successResponse(
      'Registration successful. Verification email has been queued.',
      buildUserResource(savedUser),
    );
  }

  async login(payload: LoginDto, request: Request) {
    const user = await this.userRepository.findOne({
      where: { email: payload.email.toLowerCase() },
      relations: { profile: true },
    });

    if (!user || !(await bcrypt.compare(payload.password, user.password))) {
      throw new UnprocessableEntityException({
        errors: {
          email: ['These credentials do not match our records.'],
        },
      });
    }

    return this.issueToken(user, request);
  }

  async logout(user: AuthenticatedUser) {
    if (user.jwtJti) {
      await this.userJwtSessionRepository.update(
        { jti: user.jwtJti },
        { revokedAt: new Date() },
      );
    } else {
      await this.userJwtSessionRepository.update(
        { userId: user.id, revokedAt: IsNull() },
        { revokedAt: new Date() },
      );
    }

    return successResponse('Successfully logged out.', null);
  }

  async refresh(user: AuthenticatedUser, request: Request) {
    if (user.jwtJti) {
      await this.userJwtSessionRepository.update(
        { jti: user.jwtJti },
        { revokedAt: new Date() },
      );
    } else {
      await this.userJwtSessionRepository.update(
        { userId: user.id, revokedAt: IsNull() },
        { revokedAt: new Date() },
      );
    }
    return this.issueToken(user, request);
  }

  async sendPasswordResetEmail(payload: PasswordResetEmailDto) {
    const user = await this.userRepository.findOne({
      where: { email: payload.email.toLowerCase() },
    });
    if (!user) {
      throw new UnprocessableEntityException({
        errors: { email: ['We can\'t find a user with that email address.'] },
      });
    }

    const plainToken = generateOpaqueToken(24);
    await this.passwordResetRepository.save({
      email: user.email,
      token: sha256(plainToken),
      createdAt: new Date(),
    });

    await this.mailService.sendPasswordResetEmail(user, plainToken);
    return { message: 'We have emailed your password reset link.' };
  }

  async resetPassword(payload: PasswordResetDto) {
    if (payload.password !== payload.password_confirmation) {
      throw new UnprocessableEntityException({
        errors: {
          password: ['The password confirmation does not match.'],
        },
      });
    }

    const reset = await this.passwordResetRepository.findOne({
      where: { email: payload.email.toLowerCase() },
    });

    if (!reset) {
      throw new UnprocessableEntityException({
        errors: { email: ['Invalid password reset token.'] },
      });
    }

    const expiryMinutes = this.config.get<number>('app.passwordResetExpiryMinutes', 60);
    const isExpired =
      !reset.createdAt || reset.createdAt.getTime() + expiryMinutes * 60_000 < Date.now();
    const incomingToken = sha256(payload.token);

    if (isExpired || !constantTimeEquals(reset.token, incomingToken)) {
      throw new UnprocessableEntityException({
        errors: { token: ['This password reset token is invalid.'] },
      });
    }

    const user = await this.userRepository.findOne({
      where: { email: payload.email.toLowerCase() },
    });
    if (!user) {
      throw new NotFoundException({
        errors: { email: ['User not found.'] },
      });
    }

    user.password = await bcrypt.hash(payload.password, 10);
    await this.userRepository.save(user);
    await this.passwordResetRepository.delete({ email: user.email });

    return { message: 'Your password has been reset.' };
  }

  async verifyEmail(id: number, expires: string, signature: string) {
    const user = await this.userRepository.findOne({ where: { id } });
    if (!user) {
      throw new NotFoundException();
    }

    const expiryTime = Number.parseInt(expires, 10);
    if (!Number.isFinite(expiryTime) || expiryTime < Date.now()) {
      throw new UnauthorizedException();
    }

    const payload = `${id}:${expires}`;
    const expectedSignature = signPayload(
      payload,
      this.config.get<string>('app.jwtSecret', 'change-me'),
    );

    if (!constantTimeEquals(expectedSignature, signature)) {
      throw new UnauthorizedException();
    }

    user.emailVerifiedAt ??= new Date();
    await this.userRepository.save(user);

    return {
      redirectUrl: `${this.config.get<string>('app.userDomain', 'http://localhost:5173')}/email-verified`,
    };
  }

  async resendVerificationEmail(user: User) {
    await this.mailService.sendVerificationEmail(user);
    return { message: 'Verification link sent.' };
  }

  async validateJwtPayload(payload: JwtPayload): Promise<AuthenticatedUser | null> {
    if (payload.type !== 'user') {
      return null;
    }

    const session = await this.userJwtSessionRepository.findOne({
      where: { jti: payload.jti },
    });
    if (!session || session.revokedAt || session.expiresAt.getTime() < Date.now()) {
      return null;
    }

    const user = await this.userRepository.findOne({
      where: { id: payload.sub },
      relations: { profile: true },
    });
    if (!user) {
      return null;
    }
    return Object.assign(user, { jwtJti: payload.jti });
  }

  async loginAdmin(email: string, password: string): Promise<Admin> {
    const admin = await this.adminRepository.findOne({
      where: { email: email.toLowerCase() },
    });

    if (!admin || admin.verified !== 'Y' || admin.deletedAt) {
      throw new UnauthorizedException();
    }

    const matches = await bcrypt.compare(password, admin.password);
    if (!matches) {
      throw new UnauthorizedException();
    }

    return admin;
  }

  createEmailVerificationLink(user: User): string {
    const expires = Date.now() +
      this.config.get<number>('app.emailVerificationExpiryMinutes', 4320) * 60_000;
    const payload = `${user.id}:${expires}`;
    const signature = signPayload(
      payload,
      this.config.get<string>('app.jwtSecret', 'change-me'),
    );
    return `${this.config.get<string>('app.url', 'http://localhost:3000')}/api/v1/user/email/verify/${user.id}?expires=${expires}&signature=${signature}`;
  }

  private async issueToken(user: User, request: Request) {
    const jti = generateOpaqueToken(16);
    const sessionId = randomUUID();
    const expiresAt = new Date(
      Date.now() + this.config.get<number>('app.jwtTtlSeconds', 3600) * 1000,
    );

    await this.userJwtSessionRepository.save({
      id: sessionId,
      userId: user.id,
      jti,
      expiresAt,
      revokedAt: null,
      ipAddress: request.ip ?? null,
      userAgent: request.get('user-agent') ?? null,
    });

    const token = await this.jwtService.signAsync({
      sub: user.id,
      email: user.email,
      type: 'user',
      jti,
    } satisfies JwtPayload);

    return { token };
  }

  private async ensureUniqueUser(email: string, phone: string | null) {
    const existingByEmail = await this.userRepository.findOne({
      where: { email: email.toLowerCase() },
      withDeleted: true,
    });
    if (existingByEmail) {
      throw new UnprocessableEntityException({
        errors: { email: ['The email has already been taken.'] },
      });
    }

    if (phone) {
      const existingByPhone = await this.userRepository.findOne({
        where: { phone },
        withDeleted: true,
      });
      if (existingByPhone) {
        throw new UnprocessableEntityException({
          errors: { phone: ['The phone has already been taken.'] },
        });
      }
    }
  }
}
