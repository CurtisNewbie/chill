import { Component, EventEmitter, OnInit, Output, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { PagingController } from '../common/paging';

@Component({
  template: `
  <mat-paginator #paginator [length]="pagingController.paging.total" [pageSize]="pagingController.paging.limit"
    [pageSizeOptions]="pagingController.PAGE_LIMIT_OPTIONS" aria-label="Select page">
  </mat-paginator>`,
  selector: 'app-controlled-paginator',
})
export class ControlledPaginatorComponent implements OnInit {

  @ViewChild("paginator", { static: true })
  paginator: MatPaginator = null;

  @Output("controllerReady")
  controllerEmitter = new EventEmitter<PagingController>();

  pagingController = new PagingController();

  constructor() {

  }

  ngOnInit(): void {
    this.pagingController.control(this.paginator);
    this.controllerEmitter.emit(this.pagingController);
  }

}
