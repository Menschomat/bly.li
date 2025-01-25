import { Component } from '@angular/core';
import { LoginMenuComponent } from './login-menu/login-menu.component';

@Component({
    selector: 'app-nav-bar',
    imports: [LoginMenuComponent],
    template: `
    <nav
      x-data="{ isOpen: false }"
      class="relative backdrop-blur-md m-2 rounded-xl"
    >
      <div class="px-6 py-4 flex justify-between items-center">
        <div
          class="transition-colors duration-300 transform text-gray-800 dark:text-gray-200"
        >
          <a href="#" class="font-montserrat font-bold text-2xl cursor-pointer">
            bly.li
          </a>
        </div>
        <app-login-menu></app-login-menu>
      </div>
    </nav>
  `
})
export class NavBarComponent {}
