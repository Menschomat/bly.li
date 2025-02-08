import { ApplicationConfig, provideZoneChangeDetection, inject, provideAppInitializer } from '@angular/core';
import {
  AuthConfig,
  OAuthStorage,
  provideOAuthClient,
} from 'angular-oauth2-oidc';
import {
  provideHttpClient,
  withInterceptorsFromDi,
} from '@angular/common/http';
import { BASE_PATH } from './api/shortn/';
import { ConfigService } from './services/config.service';
import { provideRouter, withDisabledInitialNavigation, withEnabledBlockingInitialNavigation, withHashLocation } from '@angular/router';
import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),

    provideRouter(routes),
    provideHttpClient(withInterceptorsFromDi()),
    provideOAuthClient({
      resourceServer: {
        allowedUrls: [window.location.origin + '/shortn'],
        sendAccessToken: true,
      },
    }),
    { provide: BASE_PATH, useValue: window.location.origin + '/shortn' },
    { provide: OAuthStorage, useFactory: storageFactory },
    provideAppInitializer(() => {
        const initializerFn = (appConfigInitializer)(inject(ConfigService));
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
