import { Component, EventEmitter, Output } from '@angular/core';
import { OAuthService } from 'angular-oauth2-oidc';

@Component({
  selector: 'app-button-primary',
  standalone: true,
  imports: [],
  template: `
    <button
      (click)="triggerBtn($event)"
      class="bg-animate p-0.5 rounded-lg from-indigo-400 via-pink-500 to-purple-500 bg-gradient-to-r"
    >
      <div
        class="px-4 py-2 bg-white dark:bg-gray-900 rounded-md text-transparent hover:bg-transparent hover:text-white"
      >
        <span
          class="bg-animate block font-montserrat font-black leading-snug bg-clip-text bg-gradient-to-r from-indigo-400 via-pink-500 to-purple-500"
          >shrink</span
        >
      </div>
    </button>
  `,
})
export class ButtonPrimaryComponent {
  @Output() click: EventEmitter<any> = new EventEmitter();

  public triggerBtn(event: Event) {
    event.stopPropagation();
    this.click.emit();
  }
}
