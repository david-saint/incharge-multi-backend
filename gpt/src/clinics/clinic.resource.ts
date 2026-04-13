import { Location } from '../database/entities/location.entity';
import { Clinic } from '../database/entities/clinic.entity';

type ClinicDistanceOptions = {
  mode: 'km' | 'mi';
  radius: number;
  actualDistance: number;
};

export function buildLocationResource(location: Location) {
  return {
    id: location.id,
    name: location.name,
    state_id: location.stateId,
    latitude: location.latitude === null ? null : Number(location.latitude),
    longitude: location.longitude === null ? null : Number(location.longitude),
    state: location.state
      ? {
          id: location.state.id,
          name: location.state.name,
          slug: location.state.slug,
        }
      : undefined,
    country: location.country
      ? {
          id: location.country.id,
          name: location.country.name,
          code: location.country.code,
        }
      : undefined,
    clinics: location.locatables?.map((locatable) => ({
      id: locatable.clinic?.id,
      name: locatable.clinic?.name,
    })),
  };
}

export function buildClinicResource(
  clinic: Clinic,
  options?: {
    includeLocations?: boolean;
    distance?: ClinicDistanceOptions;
    locationsCount?: number;
  },
) {
  return {
    id: clinic.id,
    name: clinic.name,
    address: clinic.address,
    latitude: clinic.latitude === null ? null : Number(clinic.latitude),
    longitude: clinic.longitude === null ? null : Number(clinic.longitude),
    created_at: clinic.createdAt,
    ...(options?.distance
      ? {
          mode: options.distance.mode,
          radius: options.distance.radius,
          search_radius: `${options.distance.radius}${options.distance.mode}`,
          actual_distance: Number(options.distance.actualDistance.toFixed(2)),
          distance: `${options.distance.actualDistance.toFixed(2)}${options.distance.mode}`,
        }
      : {}),
    ...(options?.includeLocations
      ? {
          locations:
            clinic.locatables
              ?.map((locatable) => locatable.location)
              .filter((location): location is Location => Boolean(location))
              .map(buildLocationResource) ?? [],
        }
      : {}),
    ...(options?.locationsCount !== undefined
      ? {
          locations_count: options.locationsCount,
        }
      : {}),
  };
}
