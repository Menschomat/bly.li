import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'shortNumber',
})
export class ShortNumberPipe implements PipeTransform {
  transform(value: number): string {
    if (value >= 1_000_000) {
      return (value / 1_000_000).toFixed(value % 1_000_000 === 0 ? 0 : 1) + 'M';
    } else if (value >= 10_000) {
      return Math.round(value / 1000) + 'k';
    } else {
      return value.toLocaleString();
    }
  }
}
