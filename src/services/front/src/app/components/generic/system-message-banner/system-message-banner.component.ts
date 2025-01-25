import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-system-message-banner',
  standalone: true,
  imports: [],
  template: ` <p>{{ dispMsg }}</p> `,
  styles: ``,
})
export class SystemMessageBannerComponent {
  @Input()
  private message: string = '';
  @Input()
  private color: 'red' | 'green' | 'yellow' | 'blue' | 'default' = 'default';

  get dispMsg() {
    return this.message;
  }
}
