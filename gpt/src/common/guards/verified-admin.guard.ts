import {
  CanActivate,
  ExecutionContext,
  Injectable,
  UnauthorizedException,
} from '@nestjs/common';
import { DataSource } from 'typeorm';
import { Admin } from '../../database/entities/admin.entity';

@Injectable()
export class VerifiedAdminGuard implements CanActivate {
  constructor(private readonly dataSource: DataSource) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const request = context.switchToHttp().getRequest();
    const adminId = request.session?.adminId;
    if (!adminId) {
      throw new UnauthorizedException();
    }

    const admin = await this.dataSource.getRepository(Admin).findOne({
      where: { id: adminId },
    });
    if (!admin || admin.verified !== 'Y' || admin.deletedAt) {
      throw new UnauthorizedException();
    }

    request.admin = admin;
    return true;
  }
}
