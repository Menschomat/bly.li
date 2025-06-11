import { Component } from '@angular/core';
import { ThemeService, ThemeMode } from '../../services/theme.service';
import { CommonModule } from '@angular/common';

const MODES: ThemeMode[] = ['system', 'lite', 'dark'];

@Component({
  selector: 'app-theme-toggle',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div
      class="flex items-center min-w-6 cursor-pointer"
      (click)="cycleMode()"
      [title]="'Mode: ' + modeLabel + ' (click to change)'"
    >
      <span class="text-xl"><i [ngClass]="modeClass"></i></span>
    </div>
  `,
})
export class ThemeToggleComponent {
  mode: ThemeMode = 'system';

  constructor(private readonly themeService: ThemeService) {
    this.themeService.mode$.subscribe((m) => (this.mode = m));
  }

  cycleMode() {
    const idx = MODES.indexOf(this.mode);
    const nextMode = MODES[(idx + 1) % MODES.length];
    this.themeService.setMode(nextMode);
  }

  get modeLabel() {
    if (this.mode === 'system') return 'System';
    if (this.mode === 'lite') return 'Light';
    return 'Dark';
  }

  get modeClass() {
    if (this.mode === 'system') return 'fa-regular fa-compass';
    if (this.mode === 'lite') return 'fa-regular fa-sun';
    return 'fa-regular fa-moon';
  }
}
