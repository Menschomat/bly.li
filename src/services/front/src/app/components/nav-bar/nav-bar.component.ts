import { Component } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { filter, map, Observable } from 'rxjs';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-nav-bar',
  standalone: true,
  imports: [CommonModule],
  template: `
    <nav
      x-data="{ isOpen: false }"
      class="relative backdrop-blur-md m-2 rounded-xl"
    >
      <div class="px-6 py-4 md:flex md:justify-between md:items-center">
        <div
          class="transition-colors duration-300 transform text-gray-800 dark:text-gray-200"
        >
          <a href="#" class="font-montserrat font-bold text-2xl cursor-pointer">
            bly.li
          </a>
        </div>
        <div
          *ngIf="!(hurz | async)"
          class="transition-colors duration-300 transform text-gray-800 dark:text-gray-200"
        >
          <a
            (click)="login()"
            href="#"
            class="text-lg cursor-pointer align-middle"
          >
            Login<i class="fa-solid fa-arrow-right-to-bracket"></i
          ></a>
        </div>
        <div
          *ngIf="hurz | async as username"
          class="transition-colors duration-300 transform text-gray-800 dark:text-gray-200"
        >
          <a class="text-lg cursor-pointer align-middle mr-5"
            >{{ username }}<i class="ml-2 fa-regular fa-user"></i
          ></a>
          <a
            (click)="logout()"
            href="#"
            class="text-lg cursor-pointer align-middle"
          >
            Logout<i class="ml-2 fa-solid fa-arrow-right-from-bracket"></i
          ></a>
        </div>
      </div>
    </nav>
  `,
})
export class NavBarComponent {
  public hurz: Observable<any>;
  public mobileIsOpen = false;
  constructor(private auth: AuthService) {
    this.hurz = auth.currentUser$.pipe(
      filter((a) => a !== null),
      map((a) => a['nickname'])
    );
  }
  public login(): void {
    this.auth.login();
  }
  public logout(): void {
    this.auth.logout();
  }
}
