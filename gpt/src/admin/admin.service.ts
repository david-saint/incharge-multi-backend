import {
  Injectable,
  NotFoundException,
  UnauthorizedException,
  UnprocessableEntityException,
} from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import * as bcrypt from 'bcrypt';
import { Repository } from 'typeorm';
import { CreateAdminDto } from '../auth/dto/create-admin.dto';
import { generateOpaqueToken } from '../common/utils/security.util';
import { buildPaginationMeta } from '../common/utils/query.util';
import { Admin } from '../database/entities/admin.entity';
import { ContraceptionReason } from '../database/entities/contraception-reason.entity';
import { EducationLevel } from '../database/entities/education-level.entity';
import { UpdateAdminDto } from './dto/update-admin.dto';

@Injectable()
export class AdminService {
  constructor(
    @InjectRepository(Admin)
    private readonly adminRepository: Repository<Admin>,
    @InjectRepository(ContraceptionReason)
    private readonly reasonRepository: Repository<ContraceptionReason>,
    @InjectRepository(EducationLevel)
    private readonly educationLevelRepository: Repository<EducationLevel>,
  ) {}

  async hasSuperAdmin() {
    const count = await this.adminRepository.count({
      where: { userType: 'Super' },
    });
    return count > 0;
  }

  async createAdmin(payload: CreateAdminDto) {
    await this.ensureUniqueEmail(payload.email);
    const admin = this.adminRepository.create({
      firstname: payload.firstname,
      lastname: payload.lastname,
      phone: payload.phone ?? null,
      email: payload.email.toLowerCase(),
      verified: payload.verified,
      userType: payload.userType,
      password: await bcrypt.hash(payload.password, 10),
      accessToken: generateOpaqueToken(24),
      rememberToken: null,
    });

    const saved = await this.adminRepository.save(admin);
    return this.serializeAdmin(saved);
  }

  async listAdmins() {
    const admins = await this.adminRepository.find({
      order: { verified: 'DESC', id: 'ASC' },
      take: 50,
    });

    return {
      data: admins.map((admin) => this.serializeAdmin(admin)),
      ...buildPaginationMeta('/allAdmins', 1, 50, admins.length),
    };
  }

  async getAdminDetails(id: number) {
    const admin = await this.adminRepository.findOne({ where: { id } });
    if (!admin) {
      throw new UnauthorizedException();
    }
    return this.serializeAdmin(admin);
  }

  async getVerifiedAdmin(id: number) {
    const admin = await this.adminRepository.findOne({ where: { id } });
    if (!admin || admin.verified !== 'Y') {
      throw new UnauthorizedException();
    }
    return admin;
  }

  async updateAdmin(id: number, payload: UpdateAdminDto) {
    const admin = await this.adminRepository.findOne({ where: { id } });
    if (!admin) {
      throw new NotFoundException();
    }

    if (payload.email && payload.email.toLowerCase() !== admin.email) {
      await this.ensureUniqueEmail(payload.email);
    }

    if (payload.firstname !== undefined) admin.firstname = payload.firstname;
    if (payload.lastname !== undefined) admin.lastname = payload.lastname;
    if (payload.phone !== undefined) admin.phone = payload.phone ?? null;
    if (payload.email !== undefined) admin.email = payload.email.toLowerCase();
    if (payload.verified !== undefined) admin.verified = payload.verified;
    if (payload.userType !== undefined) admin.userType = payload.userType;
    if (payload.accessToken !== undefined) {
      admin.accessToken = payload.accessToken;
    }
    const saved = await this.adminRepository.save(admin);
    return this.serializeAdmin(saved);
  }

  async referenceData() {
    const [reasons, educationLevels] = await Promise.all([
      this.reasonRepository.find({ order: { id: 'ASC' } }),
      this.educationLevelRepository.find({ order: { id: 'ASC' } }),
    ]);
    return {
      reasons: reasons.map((reason) => ({ id: reason.id, value: reason.value })),
      educationLevels: educationLevels.map((item) => ({ id: item.id, name: item.name })),
    };
  }

  private async ensureUniqueEmail(email: string) {
    const existing = await this.adminRepository.findOne({
      where: { email: email.toLowerCase() },
      withDeleted: true,
    });
    if (existing) {
      throw new UnprocessableEntityException({
        errors: { email: ['The email has already been taken.'] },
      });
    }
  }

  private serializeAdmin(admin: Admin) {
    return {
      id: admin.id,
      firstname: admin.firstname,
      lastname: admin.lastname,
      phone: admin.phone,
      email: admin.email,
      verified: admin.verified,
      userType: admin.userType,
      accessToken: admin.accessToken,
      created_at: admin.createdAt,
      updated_at: admin.updatedAt,
    };
  }
}
