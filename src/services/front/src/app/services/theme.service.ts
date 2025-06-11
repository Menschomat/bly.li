// src/app/services/theme.service.ts

import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export type ThemeMode = 'lite' | 'dark' | 'system';

@Injectable({
  providedIn: 'root',
})
export class ThemeService {
  private readonly modeSubject = new BehaviorSubject<ThemeMode>(
    this.loadMode()
  );
  public mode$ = this.modeSubject.asObservable();

  constructor() {
    this.applyTheme(this.modeSubject.value);

    window
      .matchMedia('(prefers-color-scheme: dark)')
      .addEventListener('change', () => {
        if (this.modeSubject.value === 'system') this.applyTheme('system');
      });
  }

  setMode(mode: ThemeMode) {
    this.modeSubject.next(mode);
    this.saveMode(mode);
    this.applyTheme(mode);
  }

  private applyTheme(mode: ThemeMode) {
    const html = document.documentElement;
    let dark = false;

    if (mode === 'system') {
      dark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    } else {
      dark = mode === 'dark';
    }

    if (dark) {
      html.classList.add('dark');
    } else {
      html.classList.remove('dark');
    }
  }

  private saveMode(mode: ThemeMode) {
    localStorage.setItem('theme-mode', mode);
  }

  private loadMode(): ThemeMode {
    return (localStorage.getItem('theme-mode') as ThemeMode) || 'system';
  }
}
