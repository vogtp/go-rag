import { Routes } from '@angular/router';
import { CollectionsListComponent } from './components/collections-list/collections-list.component';
import { SearchComponent } from './components/search/search.component';

export const routes: Routes = [
  { path: 'list', component: CollectionsListComponent },
  { path: 'query/:collection', component: SearchComponent },
  { path: '',   redirectTo: '/list', pathMatch: 'full' },
];
