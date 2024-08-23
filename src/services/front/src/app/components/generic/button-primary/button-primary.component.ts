import { Component, EventEmitter, Output } from '@angular/core';
import { OAuthService } from 'angular-oauth2-oidc';

@Component({
  selector: 'app-button-primary',
  standalone: true,
  imports: [],
  templateUrl: './button-primary.component.html',
  styleUrl: './button-primary.component.scss',
})
export class ButtonPrimaryComponent {
  @Output() click: EventEmitter<any> = new EventEmitter();

  public triggerBtn(event: Event) {
    event.stopPropagation();
    this.click.emit();
  }
}
