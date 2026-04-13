import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Algorithm } from './entities/algorithm.entity';
import { Admin } from './entities/admin.entity';
import { Clinic } from './entities/clinic.entity';
import { ContraceptionReason } from './entities/contraception-reason.entity';
import { Country } from './entities/country.entity';
import { EducationLevel } from './entities/education-level.entity';
import { Faq } from './entities/faq.entity';
import { FaqGroup } from './entities/faq-group.entity';
import { Locatable } from './entities/locatable.entity';
import { Location } from './entities/location.entity';
import { PasswordReset } from './entities/password-reset.entity';
import { Profile } from './entities/profile.entity';
import { State } from './entities/state.entity';
import { UserJwtSession } from './entities/user-jwt-session.entity';
import { User } from './entities/user.entity';
import { DatabaseSeedService } from './database.seed.service';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      User,
      Profile,
      Clinic,
      Location,
      Locatable,
      State,
      Country,
      ContraceptionReason,
      EducationLevel,
      FaqGroup,
      Faq,
      Algorithm,
      Admin,
      PasswordReset,
      UserJwtSession,
    ]),
  ],
  providers: [DatabaseSeedService],
  exports: [TypeOrmModule, DatabaseSeedService],
})
export class DatabaseModule {}
