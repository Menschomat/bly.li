import { Component } from '@angular/core';
import { NavBarComponent } from './components/nav-bar/nav-bar.component';
import { ThemeService } from './services/theme.service';
import { RouterOutlet } from '@angular/router';
import { CommonModule } from '@angular/common';
import { ThemeToggleComponent } from './components/theme-toggle/theme-toggle.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, NavBarComponent, RouterOutlet, ThemeToggleComponent],
  host: {
    class: 'flex flex-col h-full text-gray-800 dark:text-gray-200',
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
      <app-theme-toggle></app-theme-toggle>
    </div>
  `,
})
export class AppComponent {
  title = 'front';
}
