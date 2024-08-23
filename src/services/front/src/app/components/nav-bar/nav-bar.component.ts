import { Component } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { filter, map, Observable } from 'rxjs';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-nav-bar',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './nav-bar.component.html',
  styleUrl: './nav-bar.component.scss',
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
