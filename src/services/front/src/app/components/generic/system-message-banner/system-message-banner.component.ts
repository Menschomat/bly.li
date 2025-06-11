import { Component, Input, OnInit } from '@angular/core';
import { NgClass } from '@angular/common'; // Import NgClass

@Component({
  selector: 'app-system-message-banner',
  standalone: true, // Make it a standalone component
  imports: [NgClass], // Import NgClass for dynamic class binding
  template: ` <p [ngClass]="getBannerClasses()">{{ dispMsg }}</p> `,
  styles: [], // Tailwind handles styling, so no component-specific styles needed here
})
export class SystemMessageBannerComponent implements OnInit {
  @Input()
  public message: string = ''; // Made public for template access, though dispMsg handles it
  @Input()
  public color: 'red' | 'green' | 'yellow' | 'blue' | 'default' = 'default';

  private readonly colorClasses: Record<string, string> = {
    red: 'bg-red-100 text-red-700 border-red-400',
    green: 'bg-green-100 text-green-700 border-green-400',
    yellow: 'bg-yellow-100 text-yellow-700 border-yellow-400',
    blue: 'bg-blue-100 text-blue-700 border-blue-400',
    default: 'bg-gray-100 text-gray-700 border-gray-400',
  };

  private readonly commonClasses: string = 'p-4 rounded-md border'; // Common classes for all banners

  constructor() {}

  ngOnInit(): void {
    if (!this.message) {
      console.warn('SystemMessageBannerComponent: Message input is empty.');
    }
  }

  get dispMsg(): string {
    return this.message;
  }

  public getBannerClasses(): string {
    const selectedColorClasses =
      this.colorClasses[this.color] || this.colorClasses['default'];
    return `${this.commonClasses} ${selectedColorClasses}`;
  }
}
