import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import {
  CollectionSearchService,
  Document,
} from '../../services/collection-search.service';

@Component({
  selector: 'app-search-result-item',
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatCardModule,
    MatChipsModule,
    MatProgressBarModule,
  ],
  templateUrl: './search-result-item.component.html',
  styleUrl: './search-result-item.component.css',
})
export class SearchResultItemComponent {
  @Input() doc: Document | undefined;

  constructor(private searchService: CollectionSearchService) {}

  loadSummary() {
    this.searchService.summary(this.doc?.UUID!).subscribe((data) => {
      this.doc!.Summary = data.Summary;
    });
  }

  ngOnInit() {
    this.doc!.Summary = 'Loading summary...';
    this.loadSummary();
  }

  showContent: boolean = false;
  showPage: boolean = false;
  showPageButton = 'Show Page';
  showContentButton = 'Show Content';

  toggleShowContent() {
    this.showContent = !this.showContent;
    if (this.showContent) {
      this.showContentButton = 'Hide Content';
    } else {
      this.showContentButton = 'Show Content';
    }
  }
  toggleShowPage() {
    this.showPage = !this.showPage;
    if (this.showPage) {
      this.showPageButton = 'Hide Page';
    } else {
      this.showPageButton = 'Show Page';
    }
  }
}
