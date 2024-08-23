import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { UrlInputComponent } from './components/url-input/url-input.component';
import { NavBarComponent } from './components/nav-bar/nav-bar.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, UrlInputComponent, NavBarComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {
  title = 'front';

}
