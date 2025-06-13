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
    class: 'flex-1 flex flex-col text-gray-800 dark:text-gray-200',
  },
  template: `
    <div
      class="fixed top-0 inset-0 bg-[radial-gradient(#d8d8d8_1px,transparent_1px)] 
             dark:bg-[radial-gradient(#2b2b2b_1px,transparent_1px)]
             [background-size:24px_24px] -z-1"
    ></div>

    <app-nav-bar class="fixed top-0 z-15 top-0 left-0 right-0"></app-nav-bar>

    <!-- 2) flex/grow + min-h-0 + overflow-auto -->
    <main class="mt-16 flex-1 flex w-full min-h-0 overflow-auto">
      <router-outlet></router-outlet>
    </main>
    <app-theme-toggle class="fixed bottom-5 right-5 z-1"></app-theme-toggle>
    <footer
      class="flex w-full items-center justify-center p-2
             text-gray-700 dark:text-white"
    >
      <div>© Mensch0 – 2025</div>
    </footer>
  `,
})
export class AppComponent {
  title = 'front';
}
