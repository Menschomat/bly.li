import { Component } from '@angular/core';
import { ThemeService, ThemeMode } from '../../services/theme.service';
import { CommonModule } from '@angular/common';

const MODES: ThemeMode[] = ['system', 'lite', 'dark'];

@Component({
  selector: 'app-theme-toggle',
  standalone: true,
  imports: [CommonModule],
  host: { class: '' },
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
  constructor(private readonly themeService: ThemeService) {}

  cycleMode() {
    const current = this.themeService.mode();
    const idx = MODES.indexOf(current);
    const nextMode = MODES[(idx + 1) % MODES.length];
    this.themeService.setMode(nextMode);
  }

  get modeLabel() {
    const mode = this.themeService.mode();
    if (mode === 'system') return 'System';
    if (mode === 'lite') return 'Light';
    return 'Dark';
  }

  get modeClass() {
    const mode = this.themeService.mode();
    if (mode === 'system') return 'fa-regular fa-compass';
    if (mode === 'lite') return 'fa-regular fa-sun';
    return 'fa-regular fa-moon';
  }
}
