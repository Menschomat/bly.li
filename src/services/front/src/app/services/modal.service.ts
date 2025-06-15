import {
  Injectable,
  ViewContainerRef,
  ComponentRef,
  Type,
} from '@angular/core';

@Injectable({ providedIn: 'root' })
export class ModalService {
  private rootContainer?: ViewContainerRef;
  private currentModal?: ComponentRef<any>;

  setRootContainer(container: ViewContainerRef): void {
    this.rootContainer = container;
  }

  open<T>(component: Type<T>, inputs?: Record<string, any>): ComponentRef<T> {
    this.close();

    if (!this.rootContainer) {
      throw new Error('Root container not set. Call setRootContainer() first');
    }

    const componentRef = this.rootContainer.createComponent(component);

    // Set input properties
    if (inputs) {
      Object.keys(inputs).forEach((key) => {
        componentRef.setInput(key, inputs[key]);
      });
    }

    this.currentModal = componentRef;
    return componentRef;
  }

  close(): void {
    if (this.currentModal) {
      this.currentModal.destroy();
      this.currentModal = undefined;
    }
  }
}
