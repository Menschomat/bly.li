import { Component, Input } from '@angular/core';
import { BaseModalComponent } from '../base-modal/base-modal.component';
import { ModalService } from '../../../services/modal.service';
import { ButtonPrimaryComponent } from '../../generic/button-primary/button-primary.component';

@Component({
  standalone: true,
  imports: [ButtonPrimaryComponent],
  template: `
    <div class="fixed inset-0 flex items-center justify-center z-[1001] p-4">
      <div class=" w-full max-w-md" (click)="$event.stopPropagation()">
        <div
          class="p-4 rounded-3xl shadow-md 
                border border-gray-100 backdrop-blur-md 
                dark:border-gray-900"
        >
          <div
            class="border-b border-gray-200 dark:border-gray-700 flex justify-between items-center"
          >
            <h3 class="text-xl font-semibold">{{ title }}</h3>
          </div>
          <p class="">{{ message }}</p>
          <div class=" flex justify-end space-x-3">
            <app-button-primary
              [color]="'blue-emerald'"
              (click)="onCancel()"
              class=""
              >Cancel</app-button-primary
            >
            <app-button-primary
              [color]="'purple-pink'"
              (click)="onConfirm()"
              class=""
              >Delete</app-button-primary
            >
          </div>
        </div>
      </div>
    </div>
  `,
})
export class ConfirmationModalComponent extends BaseModalComponent {
  @Input() title: string = 'Confirmation';
  @Input() message: string = 'Are you sure?';
  @Input() onConfirmHandler?: () => void;

  constructor(modalService: ModalService) {
    super(modalService);
  }

  onConfirm(): void {
    this.onConfirmHandler?.();
    this.close();
  }

  onCancel(): void {
    this.close();
  }
}
