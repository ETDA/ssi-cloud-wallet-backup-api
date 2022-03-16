import {table} from "console";
import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
	return knex.schema.alterTable('vcs', (table) => {
		table.text('schema_type').notNullable().alter()
		table.text('jwt').notNullable().alter()
	})
}


export async function down(knex: Knex): Promise<void> {
	return knex.schema.alterTable('vcs', (table) => {
		table.string('schema_type').notNullable().alter()
		table.string('jwt').notNullable().alter()
	})
}

