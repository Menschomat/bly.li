import { HttpClient } from '@angular/common/http';
import { Injectable, isDevMode } from '@angular/core';
import { firstValueFrom } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ConfigService {
  private config: any;

  constructor(private http: HttpClient) {}

  async loadAppConfig(): Promise<void> {
    const configFileName = isDevMode() ? 'dev.config.json' : 'config.json';
    return firstValueFrom(this.http.get(`./assets/${configFileName}`))
      .then((data) => {
        this.config = data;
        console.log('Config Loaded');
      })
      .catch((err) => {
        console.error(`Could not load ${configFileName}`, err);
      });
  }

  getConfig() {
    return this.config;
  }
}
