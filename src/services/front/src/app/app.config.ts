import {
  APP_INITIALIZER,
  ApplicationConfig,
  provideZoneChangeDetection,
} from '@angular/core';
import {
  AuthConfig,
  OAuthStorage,
  provideOAuthClient,
} from 'angular-oauth2-oidc';
import {
  provideHttpClient,
  withInterceptorsFromDi,
} from '@angular/common/http';
import { BASE_PATH } from './core/api/v1';
import { ConfigService } from './services/config.service';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),

    //provideRouter(routes),
    provideHttpClient(withInterceptorsFromDi()),
    provideOAuthClient({
      resourceServer: {
        allowedUrls: [window.location.origin + '/shortn'],
        sendAccessToken: true,
      },
    }),
    { provide: BASE_PATH, useValue: window.location.origin + '/shortn' },
    { provide: OAuthStorage, useFactory: storageFactory },
    {
      provide: APP_INITIALIZER,
      useFactory: appConfigInitializer,
      multi: true,
      deps: [ConfigService],
    },
  ],
};

export function appConfigInitializer(appConfig: ConfigService): Function {
  return () => {
    return appConfig.loadAppConfig();
  };
}

// We need a factory, since localStorage is not available during AOT build time.
export function storageFactory(): OAuthStorage {
  return localStorage;
}
