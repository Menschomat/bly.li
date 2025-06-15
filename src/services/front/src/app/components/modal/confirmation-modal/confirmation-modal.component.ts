import { Component } from '@angular/core';
import { BaseModalComponent } from '../base-modal/base-modal.component';


@Component({
  selector: 'app-confirmation-modal',
  template: `
    <div class="modal-backdrop" (click)="close()"></div>
    
    <div class="modal" [class.active]="true">
      <div class="modal-header">
        <h2>{{ title }}</h2>
        <button class="close-button" (click)="close()">&times;</button>
      </div>
      
      <div class="modal-body">
        <p>{{ message }}</p>
        <div class="modal-footer">
          <button (click)="onCancel()">Cancel</button>
          <button (click)="onConfirm()">Confirm</button>
        </div>
      </div>
    </div>
  `,
  styles: ``,
})
export class ConfirmationModalComponent extends BaseModalComponent {
  title = 'Confirmation';
  message = 'Are you sure?';
  onConfirmHandler!: () => void;

  onConfirm(): void {
    this.onConfirmHandler?.();
    this.close();
  }

  onCancel(): void {
    this.close();
  }
}