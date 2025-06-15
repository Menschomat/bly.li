import {
  ApplicationConfig,
  provideZoneChangeDetection,
  inject,
  provideAppInitializer,
} from '@angular/core';
import {
  OAuthStorage,
  provideOAuthClient,
} from 'angular-oauth2-oidc';
import {
  provideHttpClient,
  withInterceptorsFromDi,
} from '@angular/common/http';
import { BASE_PATH } from './api/';
import { ConfigService } from './services/config.service';
import {
  provideRouter,
} from '@angular/router';
import { routes } from './app.routes';
import { provideAnimations } from '@angular/platform-browser/animations';

export const appConfig: ApplicationConfig = {
  providers: [
    provideAnimations(),
    provideZoneChangeDetection({ eventCoalescing: true }),

    provideRouter(routes),
    provideHttpClient(withInterceptorsFromDi()),
    provideOAuthClient({
      resourceServer: {
        allowedUrls: [window.location.origin],
        sendAccessToken: true,
      },
    }),
    { provide: BASE_PATH, useValue: window.location.origin },
    { provide: OAuthStorage, useFactory: storageFactory },
    provideAppInitializer(() => {
      const initializerFn = appConfigInitializer(inject(ConfigService));
      return initializerFn();
    }),
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
