import { Component } from '@angular/core';
import { AuthService } from '../../../services/auth.service';
import { filter, map, tap, Observable } from 'rxjs';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-login-menu',
  imports: [CommonModule, RouterLink, RouterLinkActive],
  template: `
    <div
      *ngIf="!(curUsrName | async)"
      class="flex flex-col sm:flex-row  items-end gap-4 transition-colors duration-300 transform text-gray-800 dark:text-gray-200"
    >
      <a (click)="login()" class="text-lg cursor-pointer align-middle">
        Login<i class="ml-2 fa-solid fa-arrow-right-to-bracket"></i
      ></a>
    </div>
    <div
      *ngIf="curUsrName | async as username"
      class=" flex flex-col sm:flex-row  items-end gap-4 transition-colors duration-300 transform text-gray-800 dark:text-gray-200"
    >
      <a
        routerLink="/dash"
        routerLinkActive="active"
        class="text-lg cursor-pointer align-middle"
      >
        <i class="ml-2 fa-solid fa-chart-line"></i
      ></a>
      <a class="text-lg align-middle "
        >{{ username }}<i class="ml-2 fa-regular fa-user"></i
      ></a>
      <a (click)="logout()" class="text-lg cursor-pointer align-middle">
        Logout<i class="ml-2 fa-solid fa-arrow-right-from-bracket"></i
      ></a>
    </div>
  `,
})
export class LoginMenuComponent {
  public curUsrName: Observable<string | undefined>;
  constructor(private readonly auth: AuthService) {
    this.curUsrName = auth.currentUser$.pipe(
      filter((a) => a !== null),
      tap((a) => console.debug(a)),
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
