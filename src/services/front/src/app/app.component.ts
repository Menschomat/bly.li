import { Component, OnInit } from '@angular/core';
import { NavBarComponent } from './components/nav-bar/nav-bar.component';
import { ThemeService, ThemeMode } from './services/theme.service';
import { RouterOutlet } from '@angular/router';
import { CommonModule } from '@angular/common';

const MODES: ThemeMode[] = ['system', 'lite', 'dark'];
@Component({
  selector: 'app-root',
  imports: [CommonModule, NavBarComponent, RouterOutlet],
  host: {
    class: 'flex flex-col h-full text-gray-800 dark:text-gray-200'
  },
  template: `
  <div
  style="z-index: -1"
  class="absolute inset-0 h-full w-full bg-white dark:bg-black bg-[radial-gradient(#d7d7d7_1px,transparent_1px)] dark:bg-[radial-gradient(#2b2b2b_1px,transparent_1px)] [background-size:24px_24px]"
></div>
    <app-nav-bar></app-nav-bar>
    <div class="flex flex-1 flex-col">
      <router-outlet></router-outlet>
    </div>
    <div class="m-2 text-gray-700 dark:text-white flex justify-between">
      <div></div>
      <div>©Mensch0 - 2025</div>
      <div
        class="flex items-center min-w-6"
        (click)="cycleMode()"
        [title]="'Mode: ' + modeLabel + ' (click to change)'"
      >
        <span class="text-xl"><i [ngClass]="modeClass"></i></span>
      </div>
    </div>
  `,
})
export class AppComponent implements OnInit {
  mode: ThemeMode = 'system';
  title = 'front';
  constructor(private readonly themeService: ThemeService) {}
  ngOnInit(): void {
    this.themeService.mode$.subscribe((m) => (this.mode = m));
  }
  public setMode(mode: ThemeMode) {
    this.themeService.setMode(mode);
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
