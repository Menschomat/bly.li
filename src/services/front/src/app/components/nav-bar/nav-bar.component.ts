import { Component } from '@angular/core';
import { LoginMenuComponent } from './login-menu/login-menu.component';

@Component({
  selector: 'app-nav-bar',
  imports: [LoginMenuComponent],
  template: `
    <nav x-data="{ isOpen: false }" class="relative  px-6 pb-4 ">
      <div
        class="flex justify-between items-center shadow-md 
                border border-gray-100 backdrop-blur-md 
                dark:border-gray-900 rounded-b-3xl z-12"
      >
        <div
          class=" 
            transition-colors 
            duration-300 
            transform 
            text-gray-800 
            dark:text-gray-200 
            p-2 px-6"
        >
          <a
            href="#"
            class="
              font-montserrat 
              font-bold 
              text-2xl 
              cursor-pointer 
              bg-animate 
              text-center 
              font-montserrat 
              font-black 
              leading-snug 
              text-transparent 
              bg-clip-text 
              bg-gradient-to-r 
              from-purple-600 
              via-pink-600  
              to-indigo-600"
          >
            bly.li
          </a>
        </div>
        <app-login-menu
          class="p-4 px-6"
        ></app-login-menu>
      </div>
    </nav>
  `,
})
export class NavBarComponent {}
