import { Component } from '@angular/core';
import { NavBarComponent } from './components/nav-bar/nav-bar.component';
import { RouterOutlet } from '@angular/router';
import { CommonModule } from '@angular/common';
import { ThemeToggleComponent } from './components/theme-toggle/theme-toggle.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, NavBarComponent, RouterOutlet, ThemeToggleComponent],
  host: {
    // 1) use h-screen to guarantee full-viewport height
    class: 'grid grid-rows-[auto_1fr_auto] h-screen text-gray-800 dark:text-gray-200',
  },
  template: `
    <div
      class="absolute inset-0 bg-[radial-gradient(#d7d7d7_1px,transparent_1px)] 
             dark:bg-[radial-gradient(#2b2b2b_1px,transparent_1px)]
             [background-size:24px_24px]">
    </div>

    <app-nav-bar class="row-start-1 row-end-2 z-10"></app-nav-bar>

    <!-- 2) flex/grow + min-h-0 + overflow-auto -->
    <main class="row-start-2 row-end-3 w-full min-h-0 overflow-auto flex">
      <router-outlet></router-outlet>
    </main>

    <footer
      class="row-start-3 row-end-4 flex items-center justify-between p-2
             text-gray-700 dark:text-white z-10">
      <div></div>
      <div>© Mensch0 – 2025</div>
      <app-theme-toggle></app-theme-toggle>
    </footer>
  `,
})
export class AppComponent {
  title = 'front';
}

