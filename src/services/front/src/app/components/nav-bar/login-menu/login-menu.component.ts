import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  NO_ERRORS_SCHEMA,
} from '@angular/core';
import { AuthService } from '../../../services/auth.service';
import { filter, map, tap, Observable } from 'rxjs';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { DropdownComponent } from '../../generic/dropdown/dropdown.component';
import { MenuItemComponent } from '../../generic/dropdown/menu-item/menu-item.component';

@Component({
  selector: 'app-login-menu',
  schemas: [CUSTOM_ELEMENTS_SCHEMA, NO_ERRORS_SCHEMA],
  imports: [
    CommonModule,
    RouterLink,
    RouterLinkActive,
    DropdownComponent,
    MenuItemComponent,
  ],
  template: `
    <div
      *ngIf="!(curUsrName | async)"
      class="flex flex-col sm:flex-row  items-end gap-4 transition-colors duration-300 transform text-gray-800 dark:text-gray-200"
    >
      <a (click)="login()" class="text-lg cursor-pointer align-middle">
        Login<i class="ml-2 fa-solid fa-arrow-right-to-bracket"></i
      ></a>
    </div>
    <app-dropdown *ngIf="curUsrName | async as username">
      <user-name>{{ username }}</user-name>
      <item-list>
        <app-menu-item routerLink="/dash" routerLinkActive="active">
          <item-content>
            <i class="mr-2 fa-solid fa-chart-line"></i>Dashboard
          </item-content>
        </app-menu-item>
        <hr class="mx-3 border-t border-gray-200 dark:border-gray-800">
        <app-menu-item (click)="logout()">
          <item-content>
            <i class="mr-2 fa-solid fa-arrow-right-from-bracket"></i>Logout
          </item-content>
        </app-menu-item>
      </item-list>
    </app-dropdown>
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
