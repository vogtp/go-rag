import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { MatIcon } from '@angular/material/icon';
import { ActivatedRoute, Router } from '@angular/router';
import { CollectionSearchResponse, CollectionSearchService } from '../../services/collection-search.service';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule, MatIcon, ReactiveFormsModule],
  templateUrl: './search.component.html',
  styleUrl: './search.component.css'
})
export class SearchComponent {

  collectionName: string = "no collection"
  searchQuery = new FormControl("")
  searchResult: CollectionSearchResponse | undefined

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private searchService: CollectionSearchService
  ) { }
  
  @Input()
  set collection(collection: string) {
    this.collectionName = collection
  }


  search() {
    let query = this.searchQuery.value!
    console.log("Query: " + query);
    this.searchService.searchCollection(this.collectionName,query).subscribe(data => {
      console.log(data);
      
      this.searchResult = data
    }
    );
  }

}
