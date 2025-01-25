import { Component } from '@angular/core';
import { AuthService } from '../../../services/auth.service';
import { filter, map, Observable } from 'rxjs';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-login-menu',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div
      *ngIf="!(curUsrName | async)"
      class="transition-colors duration-300 transform text-gray-800 dark:text-gray-200"
    >
      <a (click)="login()" href="#" class="text-lg cursor-pointer align-middle">
        Login<i class="fa-solid fa-arrow-right-to-bracket"></i
      ></a>
    </div>
    <div
      *ngIf="curUsrName | async as username"
      class=" flex flex-col sm:flex-row  items-end gap-4 transition-colors duration-300 transform text-gray-800 dark:text-gray-200"
    >
      <a class="text-lg cursor-pointer align-middle "
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
  `,
})
export class LoginMenuComponent {
  public curUsrName: Observable<string | undefined>;
  constructor(private auth: AuthService) {
    this.curUsrName = auth.currentUser$.pipe(
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
