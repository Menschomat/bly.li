// src/app/services/theme.service.ts

import {
  Injectable,
  computed,
  effect,
  signal,
} from '@angular/core';
import { toObservable } from '@angular/core/rxjs-interop';
import { Observable } from 'rxjs';

export type ThemeMode = 'lite' | 'dark' | 'system';

@Injectable({
  providedIn: 'root',
})
export class ThemeService {
  /** Current theme mode */
  readonly mode = signal<ThemeMode>(this.loadMode());
  /** Observable view of the current mode for legacy consumers */
  readonly mode$: Observable<ThemeMode> = toObservable(this.mode);

  /** Derived signal representing the active theme */
  readonly theme = computed<'dark' | 'lite'>(() => {
    const m = this.mode();
    if (m !== 'system') return m;
    return window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'lite';
  });
  /** Observable view of the current theme for legacy consumers */
  readonly theme$: Observable<'dark' | 'lite'> = toObservable(this.theme);

  constructor() {
    effect(() => this.applyTheme(this.mode()));

    window
      .matchMedia('(prefers-color-scheme: dark)')
      .addEventListener('change', () => {
        if (this.mode() === 'system') this.applyTheme('system');
      });
  }

  setMode(mode: ThemeMode) {
    this.mode.set(mode);
    this.saveMode(mode);
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
