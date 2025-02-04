import { Component } from '@angular/core';
import { NavBarComponent } from './components/nav-bar/nav-bar.component';

import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  imports: [NavBarComponent, RouterOutlet],
  host: { class: 'flex flex-col h-full text-gray-800 dark:text-gray-200' },
  template: ` <app-nav-bar></app-nav-bar>
    <div class="flex flex-1 flex-col">
      <router-outlet></router-outlet>
    </div>
    <div class="m-2 text-gray-700 dark:text-white flex justify-center">
      <div>©Mensch0 - 2025</div>
    </div>`,
})
export class AppComponent {
  title = 'front';
}
